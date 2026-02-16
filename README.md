# go-msr - library for r/w access to MSR

Continuation/rewrite of [gomsr](https://github.com/fearful-symmetry/gomsr). Since it does not seem to be actively maintained, and we actively use it in
[go-linux-lowlevel-hw](https://github.com/9elements/go-linux-lowlevel-hw), here is the rewrite. Functionality is same as in the predecessor with following changes:

- No direct usage of syscalls, standard library uses syscalls under the hood anyways [1](https://cs.opensource.google/go/go/+/master:src/internal/poll/fd_unix.go;drc=f65692ea562bf24c21ae46854e98584dd4bcc201;l=176) [2](https://cs.opensource.google/go/go/+/master:src/os/file_posix.go;l=53;drc=f65692ea562bf24c21ae46854e98584dd4bcc201?q=pwrite&sq=&ss=go%2Fgo), and it makes sure for us that these are handled safely.
- 
