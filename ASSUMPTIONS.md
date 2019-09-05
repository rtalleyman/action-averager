# action-averager Assumptions and Design Decisions

Here is all of the different assumptions and design decisions that were
made when designing this project. There are two different sections Assumptions,
which lists all of the different assumptions made and Design Decisions which
documents why certain decisions were made in regards to design.

## Assumptions

* There is an arbitrary amount of different actions that can be provided.
* The function signatures listed in the project description are explicit and
are not to be changed in regards to input arguments and returns.
* Input must match exactly the form described in the project description.
Anything that does not follow that form will be rejected.
* Times for actions can be greater than or equal to 0, but not less than 0.
* Times provided can have decimal points in them.
* If actions are not provided then an empty json array, "[]", will be returned.
* If errors are encountered with input the averager will reject that input,
but will continue averaging in regards to the previous and following correct
inputs that are provided. It is up to the user to handle errors when invalid
input is provided.

## Design Decisions

### Product

* The package exports both an averager struct and interface to allow users to
use dependency injection in their tests that use the averager and to allow
users to extend or provide new averagers that satisfy the interface.
* Pointers are used as much as possible to gain efficiency by not copying
underlying data structures.
* AddAction delays locking the datastore until after the action is processed,
to reduce the time the lock is kept locked.
* defers are used when unlocking as a future proofing method to prevent unlocks
from never being called or being called at the wrong time.
* GetStats does not unlock until after marshaling the output (even though it
could) out of future proofing reasons, since its safer to hace locks and defer unlocks
next to each other. If this caused performance issues this could be changed, but
probably GetStats will be called a lot less than AddAction so this performance
hit is probably fine.
* If the input of AddAction and GetStats was []bytes instead of string the
conversion of []bytes to strings could be removed, but as mentioned above
function signatures and returns are treated explicitly.
* float64 was used for averages, since averages should not be rounded and 64
bits gives flexibility if this was ever used for very large actions.
* All numbers were treated as float64 to cut down on type conversion costs.
* Locks were used instead of channels out of simplicity for the datastore.
* main.go is only to be used as an example of how to use this package.
* Divisions are expensive so averages are only calculated when GetStats is
called not every time an action is added.

### Tests

* A BDD test framework was used to allow for clarity and organization of test
cases.
* ginkgo and gomega were used as the BDD framework, since they are the
standard in go.
* Tests tend to be more verbose and less elegant in order to make them simpler
to read and comprehend.
* Test cases are run in parallel in order to improve performance.
* Test cases will continue to run after failure to catch other failures or
cascading failures.
* Test cases are randomized every run so there is not any assumed
preconditions before a test is run.
* Tests are run with a race detector, since there are concurrent tests.
