# action-averager Assumptions and Design Decisions

Here is all of the different assumptions and design decisions that were
made when designing this project. There are two different sections Assumptions,
which list all of the different assumptions made and Design Decisions which
documents why certain decisions were made in regards to design.

## Assumptions

* There is an arbitrary amount of different actions that can be provided.
* The function signatures listed in the project description are explicit and
are not to be changed in regards to input arguments and returns.
* Input must match exactly the form described in the project description.
* Times for actions can be greater than or equal to 0, but not less than.
* Times provided can have decimal points associated with them.
* If actions are not provided then and empty json array "[]" will be returned.
* If errors are encountered with input the averager will reject that input,
but will continue averaging in regards to the previous and following correct
inputs that are provided. It is up to the user handle errors when invalid
input is provided.

## Design Decisions

### Product

* The package exports both an averager struct and interface to allow users to
user dependency injection in their tests that use the averager and to allow
users to extend or provide new averagers that satisfy the interface.
* Pointers are used as much as possible to gain efficiency by not copying
data structures.
* AddAction delays locking the datastore until after the action is processed,
to reduce the time the lock is kept locked.
* defers are used when unlocking as a future proofing method to prevent unlocks
from never being called or being called at the wrong time.
* GetStats does not unlock until after marshaling the output (even though it
could) out of style reasons, since I like to have the locks and defer unlocks
next to eachother. If this caused performance issues this could be changed, but
I figured GetStats will be called a lot less than AddAction so this performance
hit is probably fine.
* If the input of AddAction and GetStats was []bytes instead of string the
conversion of []bytes to strings could be removed, but as mentioned above I
treated the function signatures and returns explicitly.
* float64 was used for averages, since I figured that averages should not be
rounded to be more exact.
* All numbers were treated as float64 to cut down on type conversion costs.
* Locks were used instead of channels out of simplicity.
* main.go is only to be used as an example of how to use this package.

### Tests

* BDD test framework was used to allow for clarity and organization of test cases.
* ginkgo and gomega were used as the BDD framework, since they are the standard in go.
* Tests tend to be more verbose and less elegant in order to make them simpler to comprehend.
* Test cases are run in parallel in order to improve performance.
* Test cases will continue to run after failure to catch other failures or cascading failures.
* Test cases are randomized every run so there is not any assumed preconditions.