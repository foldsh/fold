pid=
trap '[[ $pid ]] && kill "$pid" && echo -n FOLD' EXIT
sleep 999 & pid=$!
wait
