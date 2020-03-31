# TimerTest
Find a BUG in Golang Timer Reset function.

The function Reset in Timer package of Golang 1.14.1 causes the process to hang.

You can see "loops forever on sched_yield sometimes(timer related)".

Address of this bug: https://github.com/golang/go/issues/38023
