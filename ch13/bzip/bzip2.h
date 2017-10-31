#include <bzlib.h>

bz_stream *bz2alloc();
void bz2free(bz_stream *s);
int bz2compress(bz_stream *s, int action, char *in, unsigned *inlen, char *out, unsigned *outlen);
