### WIP

library capable of loading petri-nets created by https://github.com/sarahtattersall/PIPE
adding support for formal validation / simulation

BACKLOG
-------

- [ ] enhance simulation using GPU
      do matrix addition with CUDA
      https://stackoverflow.com/questions/36588411/cuda-matrix-addition
      https://github.com/barnex/cuda5

COMPLETE
--------
- [x] build library for reading & writing elementary petri-nets
- [X] add test simulation to check for boundedness in tic-tac-toe
- [x] add better testing using example xml files
- [x] import pnml as a vector state machine

ICEBOX
--------
- [ ] add random walk simulation to check for boundedness on large models
- [ ] extend to support inhibitor arcs and other petri-net enhancements
- [ ] add code generator to allow state machines to be generated from pnml and made available as golang modules
