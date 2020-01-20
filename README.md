imgx - Collection of Go modules to supplement the stdlib image packages
=======================================================================

I've accumulated a lot of little tools for dealing with images since starting
to play around with them in Go. Go's not the best language for this stuff
performance-wise (it's actually a wee bit slow for some of the heavy lifting),
but the ergonomics of the langauge seem to gel well with my overall set of
preferences for writing code, so I persist with it in spite of the gaps, and
have developed a number of helpers, tools and coping strategies (some of which
are contained herein).

Each subfolder is its own Go module. This is an artifact of GitHub's lack
of support for folder-based namespacing for repositories (which I would very
much prefer to use for this).

No guarantees whatsoever are made about the stability of _any_ APIs contained
in this repo. If you require stability, it is _strongly_ recommended that you
copy and paste the required modules into your project's `internal` folder. They
should be small and discrete enough, with minimal interdependencies.

Individual modules are separately licensed and documented (reflecting my
preference to use separate repositories), but if there is no LICENSE file
contained therein, the code is licensed under the LICENSE at the top-level.


## Probably incomplete list of modules:

- `rgba`: Alternative image format to `image.RGBA`. I feel `rgba` produces more
  readable code with less array-index gymnastics. I use this as the basis for a
  lot of further processing.
- `termpalette`: Contains `color.Palette` instances for the 256-color terminal
  palette and the 16-color terminal palette using the widely used xterm formula.
- `testimg`: Random image generation that I use a lot in testing.

