convert $1 txt: | sed 's/\([0-9]\+\),\([0-9]\+\).*\(black\|white\).*/\1;\2;\3/' | sed -r 's/black/0/' | sed -r 's/white/1/'
