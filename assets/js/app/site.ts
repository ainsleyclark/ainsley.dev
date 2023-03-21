/**
 * site.ts
 *
 * @author Ainsley Clark
 * @author URL:   https://ainsley.dev
 * @author Email: hello@ainsley.dev
 */

import { IsTouchDevice } from '../util/css';

require('./../vendor/modernizr');
import { Params } from '../params';
import { Navigation } from '../components/nav';
import { Cursor } from '../animations/cursor';
import { Skew } from '../animations/skew';
import { FitText } from '../components/fit-text';
import { Card } from '../components/card';
import { Arrow } from '../animations/arrow';
import { Collapse, CollapseOptions } from '../components/accordion';
import { Log } from '../util/log';
import { beforeAfter } from '../components/before-after';
import { bookmark } from '../components/bookmark';
import { buttonGoBack } from '../components/button';
import { copyToClipboard } from '../components/copy';
import { lazyImages } from '../components/image';
import { video } from '../components/video';
import { WebVitals } from '../analytics/web-vitals';
import { Animations } from '../animations/text';
import { Elements } from '../util/els';
import { Barba } from './barba';
import Scroll from './scroll';
import { ITransitionData } from '@barba/core';
import { OnScrollEvent } from 'locomotive-scroll';

/**
 *
 */
declare global {
	interface Window {
		plausible: (args: string) => unknown;
	}
}

type ThemeColour = 'black' | 'white';

/**
 * App is the main type for the site which bootstraps the
 * application and initialises types.
 */
class App {
	/**
	 * Determines if the app has already been booted.
	 *
	 * @private
	 */
	private hooksAdded = false;

	/**
	 * The page transition type.
	 *
	 * @private
	 */
	private barba: Barba;

	/**
	 * The navigation type used for determining if the nav
	 * is open and playing animations.
	 *
	 * @private
	 */
	private nav: Navigation;

	/**
	 * The cursor type used for destroy and initialising
	 * cursor animations.
	 *
	 * @private
	 */
	private cursor: Cursor;

	/**
	 * Bootstraps the application.
	 */
	boot(): void {
		Log.debug('Booting ainsley.dev JS');
		Log.debug('JS Params', Params);

		// Instantiate JS
		this.removeJSClasses();

		// Hooks
		if (!this.hooksAdded) {
			this.once();
		}

		// Classes
		this.cursor = new Cursor();
		new Skew();
		new FitText();
		new Card();
		new Arrow();
		new Collapse({
			accordion: true,
			container: '.accordion',
			item: '.accordion-item',
			inner: '.accordion-content',
			activeClass: 'accordion-item-active',
		} as CollapseOptions);

		this.preventInternalLinks();

		// Functions
		beforeAfter();
		bookmark();
		buttonGoBack();
		copyToClipboard();
		lazyImages();
		video();

		// Animations
		this.initAnimations();

		// Analytics
		this.webVitals();
	}

	/**
	 * Initialises types only once (singletons).
	 *
	 * @private
	 */
	private once(): void {
		this.barba = new Barba();
		this.nav = new Navigation();
		this.barba.init(this.nav);
		this.before();
		this.beforeEnter();
		this.after();
		this.mouseMoveHandler();
		this.hooksAdded = true;
	}

	/**
	 * Removes/adds the no Javascript classes from
	 * the HTML element.
	 *
	 * @private
	 */
	private removeJSClasses(): void {
		Elements.HTML.classList.remove('no-js');
		Elements.HTML.classList.add('js');
	}

	/**
	 * Initialises the main pages animations. If the navigational
	 * element is currently animating, a delay will be applied.
	 *
	 * @private
	 */
	private initAnimations(): void {
		const animations = new Animations(),
			timeout = this.nav.duration() / 2 - 350;
		setTimeout(() => animations.play(), this.nav.isAnimating ? timeout : 0);
	}

	/**
	 * Before Enter Hook - Before enter transition/view
	 *
	 * @private
	 * @see https://barba.js.org/docs/advanced/hooks/
	 */
	private beforeEnter(): void {
		this.barba.hooks.beforeEnter((data: ITransitionData) => {
			Scroll.destroy();
			this.reloadJS(data.next.container);
			if (!this.hasSmoothScroll()) {
				Elements.HTML.style.scrollBehavior = "smooth";
			}
		});
	}

	/**
	 * Before Hook - Before everything
	 *
	 * @private
	 * @see https://barba.js.org/docs/advanced/hooks/
	 */
	private before(): void {
		this.barba.hooks.before(() => {
			if (this.nav.isOpen) {
				this.nav.play();
			}
			if (!this.hasSmoothScroll()) {
				Elements.HTML.style.scrollBehavior = "initial";
			}
			this.cursor.destroy();
		});
	}

	/**
	 * After Hook - After everything
	 *
	 * @private
	 * @see https://barba.js.org/docs/advanced/hooks/
	 */
	private after(): void {
		this.barba.hooks.after((data: ITransitionData) => {
			Elements.HTML.scrollTop = 0;
			Elements.Body.scrollTop = 0;
			data.next.container.scrollTop = 0;
			Scroll.init(data.next.container);
			this.triggerPageView();
			this.boot();
		});
	}

	/**
	 *
	 * @private
	 */
	private webVitals(): void {
		WebVitals({
			enable: Params.isProduction,
			analyticsId: Params.vercelAnalyticsID,
			debug: Params.appDebug,
		});
	}

	/**
	 * Prevents reloading of the page if the link clicked
	 * is the current link, to avoid the page
	 * reloading.
	 * TODO: Check if we can do this in Barba.
	 *
	 * @private
	 */
	private preventInternalLinks(): void {
		document.querySelectorAll<HTMLAnchorElement>('a[href]').forEach((link) => {
			link.addEventListener('click', (e: Event) => {
				const link = e.currentTarget as HTMLAnchorElement;
				if (link.href === window.location.href && this.nav.isOpen) {
					e.preventDefault();
					e.stopPropagation();
					this.nav.play();
				}
			});
		});
	}

	/**
	 * Finds all scripts within the next container and
	 * appends them to ensure JS is loaded after
	 * a page transition.
	 *
	 * @param container
	 * @private
	 */
	private reloadJS(container: HTMLElement): void {
		const js = container.querySelectorAll('script');
		js.forEach((item: HTMLScriptElement) => {
			if (item.src.includes('app')) {
				return;
			}
			const script = document.createElement('script');
			script.src = item.src;
			container.appendChild(script);
		});
	}

	/**
	 * Triggers a page view dynamically with Plausible.
	 *
	 * @private
	 */
	private triggerPageView(): void {
		if (typeof window.plausible === 'function') {
			Log.debug('Triggering Plausible page-view');
			window.plausible('pageview');
		}
	}

	/**
	 * Returns the current page theme colour.
	 *
	 * @private
	 */
	private getThemeColour(container: HTMLElement): ThemeColour {
		return <'black' | 'white'>container.querySelector('main').getAttribute('data-theme') ?? 'black';
	}

	/**
	 * Updates the client mouse positions to the html element.
	 *
	 * @private
	 */
	private mouseMoveHandler(): void {
		const html = Elements.HTML;
		document.addEventListener('mousemove', (e) => {
			html.setAttribute('client-x', e.clientX.toString());
			html.setAttribute('client-y', e.clientY.toString());
			html.setAttribute('x', e.x.toString());
			html.setAttribute('y', e.y.toString());
			html.setAttribute('percent-x', ((e.x / window.innerWidth) * 100).toString());
			html.setAttribute('percent-y', ((e.y / window.innerHeight) * 100).toString());
		});
	}

	/**
	 * Determines if smooth scroll is enabled.
	 *
	 * @private
	 */
	private hasSmoothScroll(): boolean {
		return !IsTouchDevice();
		// return Elements.HTML.classList.contains("has-smooth-scroll")
	}
}

export default new App();
