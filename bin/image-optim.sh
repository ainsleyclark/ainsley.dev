#!/bin/bash
#
# Shell script to convert images in the public folder
# to WebP and optimise JPG's and PNG's.
#
# ainsley.dev - 16/09/2021

# Set Variables
PUBLIC_PATH="./public"
MAXIMUM_JPG_SIZE=250
WEBP_QUALITY=80
PNG_OPTIMIZATION_LEVEL=2

# Convert to WebP.
echo '--------------------------------------------'
echo 'Processing files and converting to webp'
echo '--------------------------------------------'
find ${PUBLIC_PATH} -type f \( -name "*.png" -or -name "*.jpg" -or -name "*.jpeg"  \) | xargs -P ${WEBP_QUALITY} -I {} sh -c 'cwebp -q 80 $1 -o "${1%.*}.webp"' _ {} \;

echo '--------------------------------------------'
echo 'Converting to AVIF'
echo '--------------------------------------------'
if hash avifenc 2>/dev/null; then
	find ${PUBLIC_PATH} -type f \( -name "*.png" -or -name "*.jpg" -or -name "*.jpeg" \) | xargs -P 8 -I {} sh -c 'avifenc --min 0 --max 63 --speed 6 -a end-usage=q -a cq-level=18 -a tune=ssim $1 "${1%.*}.avif"' _ {} \;
else
	echo "Install avifenc to convert to images to .avif"
fi

# Compress JPGS's.
echo '--------------------------------------------'
echo 'Compressing JPG images'
echo '--------------------------------------------'
if hash jpegoptim 2>/dev/null; then
	for image in $(find ${PUBLIC_PATH} -type f \( -name "*.jpg" -or -name "*.jpeg"  \)); do
		# Remove all metadata and try to optimize jpeg image to match the maximum size.
		jpegoptim --strip-all --size=$MAXIMUM_JPG_SIZE $image
	done;
else
	echo "Install jpegoptim to optimize JPEG images"
fi

# Compress PNG's.
echo '--------------------------------------------'
echo 'Compressing PNG images'
echo '--------------------------------------------'
if hash optipng 2>/dev/null; then
	for image in $(find ${PUBLIC_PATH} -type f \( -name "*.png" \)); do
		# Optimize PNG with a give level (higher = slower) and remove all metadata.
		optipng -clobber -strip all -o $PNG_OPTIMIZATION_LEVEL $image
	done;
else
	echo "Install optipng to optimize PNG images"
fi

echo '--------------------------------------------'
echo 'Optimising SVG images'
echo '--------------------------------------------'
if hash svgo 2>/dev/null; then
	svgo -f ${PUBLIC_PATH}
else
	echo "Install svgo to optimize SVG images"
fi
