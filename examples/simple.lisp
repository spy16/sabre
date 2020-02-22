; define a value binding
(def π 3.1412)

; define a lambda function
(def echo (fn* [msg] msg))

; index access for a vector
([1 2 3] 0)

(def hello (fn*
    ([] "Hello")
    ([name] (let* [user-name name
                  msg (str "Hello " user-name "!")]
                msg))
))



(def defn (fn* defn [name args & body]
    (do
        (if (not (symbol? name))
            (throw "name must be symbol, not " ))
        (if (not (vector? args))
            (throw "args must be a vector, not "))
        `(def ~name (fn* ~args ~body)))))
