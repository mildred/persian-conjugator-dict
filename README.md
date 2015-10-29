Persian Verbs Downloader
========================

Download verb conjugaison from
[http://www.jahanshiri.ir/pvc/](http://www.jahanshiri.ir/pvc/) and construct a
StarDict dictionary from it.


Simple Usage
------------

Build the code:

    go build .

Download the words:

    ./persian-conjugator-dict

Run `stardict-editor` to convert the generated `.txt` file to a StarDict
dictionary. Because the definition is HTML, you'll have to update the generated
`.ifo` file to contain:

    sametypesequence=h

Instead of:

    sametypesequence=m

Advanced Usage
--------------

You can change multiple settings using command line options, just run the
program with `--help` to see what you can do.

