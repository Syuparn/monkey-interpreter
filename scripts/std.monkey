"from ''Writing An Interpreter In Go'' "
let map = fn(arr, f) {
    let iter = fn(arr, acc) {
        if (len(arr) == 0) {
            acc;
        } else {
            iter(rest(arr), push(acc, f(first(arr))));
        };
    };
    iter(arr, []);
};

"from ''Writing An Interpreter In Go'' "
let reduce = fn(arr, initial, f) {
    let iter = fn(arr, result) {
        if (len(arr) == 0) {
            result;
        } else {
            iter(rest(arr), f(result, first(arr)))
        };
    };
    iter(arr, initial);
};

"from ''Writing An Interpreter In Go'' "
let sum = fn(arr) {
    reduce(arr, 0, fn(init, el) { init + el });
};

let filter = fn(arr, cond) {
    let iter = fn(arr, result) {
        if (len(arr) == 0) {
            result;
        } else {
            iter(
                rest(arr),
                if (cond(first(arr))) {
                    push(result, first(arr));
                } else {
                    result;
                }
            );
        };
    };
    iter(arr, []);
};

let extend = fn(arrOne, arrTwo) {
    let iter = fn(arr, extended) {
        if (len(arr)==0) {
            extended;
        } else {
            iter(rest(arr), push(extended, first(arr)));
        }
    };
    iter(arrTwo, arrOne);
};

"NOTE: if (!0) {} == null"
let compactmap = fn(arr, f) {
    filter(map(arr, f), fn(x) { x != if (!0) {} });
}

let flatmap = fn(arr, f) {
    flatten(map(arr, f));
};

let flatten = fn(arr) {
    let iter = fn(arr, flat) {
        if (type(arr) != "ARRAY") {
            return iter([], extend(flat, [arr]));
        }

        if (len(arr) == 0) {
            flat;
        } else {
            iter(rest(arr), extend(flat, flatten(first(arr))));
        };
    };
    iter(arr, []);
};

let abs = fn(x) {
    if (x > 0) {
        x;
    } else {
        -x;
    };
};

let zip = fn(arrOne, arrTwo) {
    let iter = fn(arrOne, arrTwo, zipped) {
        if (len(arrOne) == 0 || len(arrTwo) == 0) {
            zipped;
        } else {
            iter(
                rest(arrOne), rest(arrTwo),
                push(zipped, [first(arrOne), first(arrTwo)])
            );
        }
    };
    iter(arrOne, arrTwo, []);
};

let enumerate = fn(arr) {
    let iter = fn(arr, i, enumerated) {
        if (len(arr) == 0) {
            enumerated;
        } else {
            iter(rest(arr), i + 1, push(enumerated, [i, first(arr)]));
        }
    };
    iter(arr, 0, []);
};

let count = fn(arr, cond) {
    len(filter(arr, cond));
};

let all = fn(arr, cond) {
    count(arr, cond) == len(arr);
};

let any = fn(arr, cond) {
    count(arr, cond) > 0;
};

let repeat = fn(arr, n) {
    let iter = fn(repeated, i) {
        if (i <= 0) {
            repeated;
        } else {
            iter(extend(repeated, arr), i - 1);
        }
    };
    iter([], n);
};
