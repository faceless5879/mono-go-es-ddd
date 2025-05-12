# My practice repo

This is just a sample, practice level so the business level complexity is very low.

That why I run transaction on repository layer. If the complexity is higher and a single execution spans on multiple aggregates, should implement tx on command layer or if the execution can be reused on multiple case then it should be domain service layer that handle the tx