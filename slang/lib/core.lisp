(ns 'core)

(def fn
    (macro* fn
        ([& decl] (cons 'fn* decl))))

(def defn (macro* defn [name & fdecl]
    (let* [func (cons 'fn (cons name fdecl))]
    `(def ~name ~func))))

(def defmacro (macro* defmacro [name & mdecl]
    (let* [macro (cons 'macro* (cons name mdecl))]
    `(def ~name ~macro))))

(defn nil? [arg] (= nil arg))

(defn empty? [coll]
    (if (nil? coll)
        true
        (nil? (first coll))))

; Type check functions
(defn same-type? [a b] (= (type a) (type b)))
(defn seq? [arg] (impl? arg types/Seq))
(defn set? [arg] (same-type? #{} arg))
(defn list? [arg] (same-type? () arg))
(defn vector? [arg] (same-type? [] arg))
(defn int? [arg] (same-type? 0 arg))
(defn float? [arg] (same-type? 0.0 arg))
(defn boolean? [arg] (same-type? true arg))
(defn string? [arg] (same-type? "" arg))
(defn keyword? [arg] (same-type? :specimen arg))
(defn symbol? [arg] (same-type? 'specimen arg))

; Type initialization functions
(defn set [coll] (apply-seq (type #{}) coll))
(defn list [& coll] (apply-seq (type ()) coll))
(defn vector [& coll] (apply-seq (type []) coll))
(defn int [arg] (to-type arg (type 0)))
(defn float [arg] (to-type arg (type 0.0)))
(defn boolean [arg] (true? arg))

(defn true? [arg]
    (if (nil? arg)
        false
        (if (boolean? arg)
            arg
            true)))

(defn not [arg] (= false (true? arg)))

(defn last [coll]
    (let* [v   (first coll)
           rem (next coll)]
        (if (nil? rem)
            v
            (last (next coll)))))
