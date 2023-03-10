/**
 * app.js
 * All custom JS for application stored here.
 *
 * @author Ainsley Clark
 * @author URL:   https://www.ainsleyclark.com
 * @author Email: info@ainsleyclark.com
 */

import { Cursor } from './animations/cursor';
import { Skew } from './animations/skew';
import { FitText } from './components/fit-text';
import { Collapse, CollapseOptions } from './components/accordion';
import { Card } from './components/card';
import { Navigation } from './components/nav';
import { Log } from './util/log';
import { Toast } from './animations/toast';
import { Arrow } from './animations/arrow';
import LocomotiveScroll from 'locomotive-scroll';
import AOS from 'aos';
import smoothscroll from 'smoothscroll-polyfill';

require('./animations/text');

/**
 * Vars
 */
const html = document.documentElement,
	body = document.body;

/*
 * Remove No JS Body Class
 */
html.classList.remove('no-js');
html.classList.add('js');

/**
 * Initialise components & types.
 */
document.addEventListener('DOMContentLoaded', () => {
	new Navigation();
	new Cursor();
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
});

/**
 * Kick off Smooth Scroll Polyfill
 */
smoothscroll.polyfill();

/**
 * Animate on Scroll
 */
AOS.init({
	// Global settings:
	disable: false, // accepts following values: 'phone', 'tablet', 'mobile', boolean, expression or function
	startEvent: 'DOMContentLoaded', // name of the event dispatched on the document, that AOS should initialize on
	initClassName: 'aos-init', // class applied after initialization
	animatedClassName: 'aos-animate', // class applied on animation
	useClassNames: false, // if true, will add content of `data-aos` as classes on scroll
	disableMutationObserver: false, // disables automatic mutations' detections (advanced)
	debounceDelay: 50, // the delay on debounce used while resizing window (advanced)
	throttleDelay: 99, // the delay on throttle used while scrolling the page (advanced)
	// Settings that can be overridden on per-element basis, by `data-aos-*` attributes:
	offset: 120, // offset (in px) from the original trigger point
	delay: 0, // values from 0 to 3000, with step 50ms
	duration: 800, // values from 0 to 3000, with step 50ms
	easing: 'ease', // default easing for AOS animations
	once: false, // whether animation should happen only once - while scrolling down
	mirror: false, // whether elements should animate out while scrolling past them
	anchorPlacement: 'top-bottom', // defines which position of the element regarding to window should trigger the animation
});

import Ukiyo from 'ukiyojs';
// CSS selector
new Ukiyo('.para', {
	scale: 1, // 1~2 is recommended
	speed: 2, // 1~2 is recommended
	willChange: true, // This may not be valid in all cases
	wrapperClass: 'para-wrapper',
	externalRAF: false,
});

/**
 * Locomotive Scroll
 */
// const scroll = new LocomotiveScroll({
// 	el: document.querySelector('[data-scroll-container]'),
// 	smooth: true,
// 	smoothMobile: false,
// 	offset: [0, 0],
// 	lerp: 0.1,
// 	touchMultiplier: 0,
// });

/**
 * Scroll
 */
const scrollAmount = 500;
body.addEventListener('scroll', () => {
	const y = body.scrollTop;
	html.style.setProperty('--scroll-y', y.toString());
	if (y > scrollAmount) {
		Log.info('Scrolled passed point');
		html.classList.add('scrolled');
		return;
	}
	html.classList.remove('scrolled');
});

/**
 * Lazy Images
 */
document.querySelectorAll('.lazy-animate').forEach((lazy) => {
	lazy.addEventListener('load', () => {
		AOS.refresh();
		lazy.classList.add('lazy-loaded');
	});
});

/**
 * Videos
 */
document.querySelectorAll('video').forEach((vid) => {
	vid.addEventListener('play', () => {
		vid.classList.add('video-playing');
	});
});

/**
 * Button - Go Back
 */
document.querySelectorAll('[data-go-back]').forEach((btn) => {
	btn.addEventListener('click', (e) => {
		e.preventDefault();
		history.back();
	});
});

/**
 * Copy to Clipboard
 */
document.querySelectorAll('[data-clipboard]').forEach((clip) => {
	const text = clip.getAttribute('data-clipboard-text');
	if (!text || text === '') {
		Log.error('Clipboard text is empty [data-clipboard-text] for el: ' + clip);
		return;
	}
	clip.addEventListener('click', () => {
		navigator.clipboard
			.writeText(text)
			.then(() => {
				Log.info('Copied text to clipboard');
				const message = clip.getAttribute('data-clipboard-message') ?? 'Copied text to clipboard';
				Toast(message);
			})
			.catch((err) => {
				Log.error('Failed to copy to clipboard: ' + err);
			});
	});
});

/**
 * Bookmark
 */
document.querySelectorAll('[data-bookmark]').forEach((bookmark) => {
	bookmark.addEventListener('click', () => {
		const userAgent = navigator.userAgent.toLowerCase();
		Toast('Press ' + (userAgent.indexOf('mac') != -1 ? 'Cmd' : 'Ctrl') + '+D to bookmark this page.');
	});
});

/**
 * Before / After
 */
document.querySelectorAll('.before-after').forEach((el) => {
	const slider = el.querySelector('.before-after-slider') as HTMLInputElement;
	if (!slider) {
		Log.error('No foreground found for before/after element');
		return;
	}
	slider.addEventListener('input', (e) => {
		const foreground = <HTMLElement>el.querySelector('.before-after-background'),
			thumb = <HTMLButtonElement>el.querySelector('.before-after-thumb');
		if (!foreground || !thumb) {
			Log.error('Element missing from before/after element');
			return;
		}
		const value = (e.target as HTMLInputElement).value;
		foreground.style.width = `${value}%`;
		thumb.style.left = `${value}%`;
	});
});
