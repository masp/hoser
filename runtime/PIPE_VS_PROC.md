Hoser is a hybrid "glue" language like bash/shellscript that has two core classes of subroutines:
- Processes
- Pipelines

A process is like a normal operating system process. Its key characteristics are:
- Inputs/outputs
- Isolated

A pipeline is a textual description of stream processing. It connects streams of data with processes to transform
those streams. Key characteristics of a pipeline are:
- DAG, not sequential (no instruction pointer)

# Relationship between processes and pipelines

### pipe -> pipe
A pipe calling another pipe is a way to re-use a smaller pipeline in a larger program. For example, imagine you had a pipe that filters explicit words
from a text stream:

```badwords.hos
module badwords;

pipe Filter(in: text) (out: text, found: lines) {
    out, bad = txt.Remove(in, "bad[0-9]+")
    out, super_bad = txt.Remove(out, "super_bad")
    ...
    found = bad + super_bad
}
```

```main.hos
module main;

pipe main(stdin: text) (stdout: text) {
    clean, bad_words = badwords.Filter(stdin)
    io.Dump("bad_words.txt", bad_words)
    stdout = clean
}
```

### pipe -> proc
If a pipe wants to perform a computation on a stream, it can call a process like `wc` to process the stream and return a new stream (map)
or a single value (reduction). The below example uses two well known processes `grep` and `wc` as examples.

```main.hos
module main;

pipe main(stdin: text) (stdout: text) {
    stdout = wc.CountLines(grep.Filter(stdin, "^ab"))
}

proc Filter(in: text)

proc CountLines(in: text)
```

### proc -> pipe
A process can itself start a pipeline. It can start it synchronously or asynchronously.

```main.hos
module main;

proc main(stdin: text) {
    result = exec Filter(stdin)

    fork Filter(stdin)
}

module badwords
pipe Filter(in: text) (out: int)

```
proc -> proc # Normal function calls in a standard programming language