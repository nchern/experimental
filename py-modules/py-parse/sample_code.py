
def foo(a, b):
    a = a + b
    print(a, b)


def bar(x):
    with open('f') as f:
        f.write('abc')
    x += 1
    print(x+1)
    x = x + 1
    return x


def nested_funcs():
    a = 0

    def _inner(a):
        return a


def multi(
    a,
    b,
        c=12):
    pass
    if a is None:
        print(123)
    if a == b and \
            True:
        print(bla)


class A:

    def bar(self, aa):
        print(aa)

    def classmethod(*, a, b, c, f):
        pass
