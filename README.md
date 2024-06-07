# Masonry Layout Image Generator

This Go program creates a masonry-style layout of images, where images fill up the width of the grid dynamically. The images are sorted by width in descending order to minimize gaps.

## Features

- Resize images to a fixed height while maintaining aspect ratio.
- Dynamically calculate the grid width based on the images in the first row.
- Sort images by width to optimize layout and minimize gaps.
- Output the final image as a PNG with a transparent background.

## Installation

### Prerequisites

- Go (version 1.16 or later)
- Git

### Steps

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/masonry-layout.git
    cd masonry-layout
    ```

2. Install the necessary Go dependencies:

    ```sh
    go get github.com/nfnt/resize
    ```

3. Build the binary:

    ```sh
    go install
    ```

## Usage

Run the binary with the desired parameters:


	masonry-layout -input /path/to/images -output output.png -maxwidth 820 -height 200


### Command-line Flags:

- input: Directory containing the input images (default is the current directory).
- output: Output image file (default is output.png).
- maxwidth: Maximum width of the output image (default is 820 pixels).
- height: Height of the thumbnails (default is 200 pixels).

## Example

To generate a masonry layout with images from the images directory, with a maximum width of 820 pixels and a thumbnail height of 200 pixels, run:

	masonry-layout -input ./images -output output.png -maxwidth 820 -height 200