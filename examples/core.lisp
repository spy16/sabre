(ns 'core)

(def nil? (fn* nil? [arg]
    (= nil arg)))

(def empty? (fn* empty? [coll]
    (if (nil? coll)
        true
        (nil? (first coll)))))

(def true? (fn* true? [arg]
    (if (nil? arg)
        false
        (if (boolean? arg)
            arg
            true))))

(def not (fn* not [arg]
    (= false (true? arg))))

(def same-type? (fn* same-type? [a b] (= (type a) (type b))))

; Type check functions
(def seq? (fn* seq? [arg] (impl? arg types/Seq)))
(def set? (fn* set? [s] (same-type? #{} s)))
(def list? (fn* list? [s] (same-type? () s)))
(def vector? (fn* vector? [s] (same-type? [] s)))
(def int? (fn* int? [arg] (same-type? 0 arg)))
(def float? (fn* float? [arg] (same-type? 0.0 arg)))
(def boolean? (fn* boolean? [arg] (same-type? true arg)))
(def string? (fn* string? [arg] (same-type? "" arg)))
(def keyword? (fn* keyword? [arg] (same-type? :keyword arg)))
(def symbol? (fn* symbol? [arg] (same-type? 'sample arg)))

; Type initialization functions
(def set (fn* set [s] (apply-seq (type #{}) s)))
(def list (fn* list [& args] (realize args)))
(def vector (fn* list [& args] (realize args)))

(def int (fn* int [arg] (to-type arg (type 0))))
(def float (fn* float [arg] (to-type arg (type 0.0))))
(def boolean (fn* boolean [arg] (true? arg)))

(def defn (fn* defn [name args & body]
    `(def ~name (fn* ~args (do (quote ~body))))))
