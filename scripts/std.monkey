"from ''Writing An Interpreter In Go'' "
let map = fn(arr, f) {
    let iter = fn(arr, acc) {
        if (len(arr)==0) {
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
        if (len(arr)==0) {
            result;
        } else {
            iter(rest(arr), f(result, first(arr)))
        };
    };
    iter(arr, initial);
};

"from ''Writing An Interpreter In Go'' "
let sum = fn(arr) {
    reduce(arr, 0, fn(init, el) {init + el});
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
    let iter = fn(arr,extended) {
        if (len(arr)==0) {
            extended;
        } else {
            iter(rest(arr), push(extended, first(arr)));
        }
    };
    iter(arrTwo, arrOne);
};

"TODO: flatmap, compactmap"