; basic functions
(assert (= 3 (eval '(int (+ 1 2)))))
(assert (true? true))
(assert (true? []))
(assert (not (true? false)))
(assert (not (true? nil)))

; sequence functions
(assert (seq? []))
(assert (not (seq? nil)))
(assert (= 1 (first [1 2 3 4])))
(assert (= [2 3 4] (next [1 2 3 4])))
(assert (= nil (next [])))
(assert (= [1 2 3 4] (cons 1 [2 3 4])))
(assert (= [1 2 3 4] (conj [1 2] 3 4)))

; threading macros
(assert (= (-> 1 (cons [2 3 4])) [1 2 3 4]))
(assert (= (->>  1 (conj [2 3 4])) [2 3 4 1]))

; basic math operators
(assert (= 3.00000 (+ 1 2)))
(assert (= 0.00000 (+)))
(assert (= 3.00000 (- 5 2)))
(assert (= -5.00000 (- 5)))
(assert (= 10.00000 (* 5 2)))
(assert (= 5.00000 (/ 10 2)))
(assert (= 0.50000 (/ 2)))
(assert (> 10 9 8 7 6 1 -1 -10))
(assert (< -10 1 2 3 4 10 23.32423432 100000))
(assert (>= 10 10 10 9 8 7 7 7 5))
(assert (<= -1.5 -1 0 0  0  0 0 0 1 2 3 4 5))

; type initialization functions
(assert (= #{} (set [])))
(assert (= #{1 2 3} (set [1 1 2 2 3])))
(assert (= [] (vector)))
(assert (= [1 2 ["hello"]] (vector 1 2 ["hello"])))
(assert (= () (list)))
(assert (= '(1 [] ["hello"] "hello") (list 1 [] ["hello"] "hello")))
(assert (= "" (str nil)))
(assert (= "" (str)))
(assert (= "1" (str 1)))
(assert (= "hello-bob" (str "hello-" "bob")))
(assert (= 1 (int 1.5677)))
(assert (= 3.00000 (float 3)))

; type checking functions
(assert (int? 10))
(assert (not (int? 10.0)))
(assert (string? ""))
(assert (not (string? nil)))
(assert (boolean? true))
(assert (boolean? false))
(assert (not (boolean? nil)))
(assert (vector? []))
(assert (not (vector? nil)))
(assert (= (type [])(type [1 2 3])))
(assert (symbol? 'hello))
