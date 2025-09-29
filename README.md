saddle4
========

A super simple auto-code formatter for my pony code.

Mostly all it does is align indents in code to the cloest 4-space stop.

The original file will be renamed to have .prev as a suffix.

The new file, with space indents fixed up, will replace the original file.

It is imperfect, but generally once you've fixed the few
things left manually, they stay fixed from then on. The 'else'
clauses are the most egregiously wrong on occassion.

installation: `go install github.com/glycerine/saddle4@latest`

---
Author: Jason E. Aten, Ph.D.

LICENSE: MIT
