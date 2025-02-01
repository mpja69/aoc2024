lines = [l.split() for l in open('input.txt') if '->' in l]


def r(out, op):
    for a, x, b, _, _ in lines:
        if op==x and out in (a,b):
            return True
    return False

def normal(a, op, b, out):
    return op == "XOR" and all(d[0] not in 'xyz' for d in (a, b, out)) 

def and_not_first(a, op, b, out):
    return op == "AND" and not "x00" in (a, b) and r(out, 'XOR')

def xor_not_first(a, op, b, out):
    return op == "XOR" and not "x00" in (a, b) and r(out, 'OR')

def last(op, out):
    return op != "XOR" and out[0] == 'z' and out != "z45"

res = []
for a, op, b, _, out in lines:
    if normal(a,op,b,out) or and_not_first(a,op,b,out) or xor_not_first(a,op,b,out) or last(op,out):
        res.append(out)


# res = [c for a, x, b, _, c in lines if
#     x == "XOR" and all(d[0] not in 'xyz' for d in (a, b, c)) or
#     x == "AND" and not "x00" in (a, b) and rr(c, 'XOR') or
#     x == "XOR" and not "x00" in (a, b) and rr(c, 'OR') or
#     x != "XOR" and c[0] == 'z' and c != "z45"]

print(*sorted(res))

# r = lambda c, y: any(y == x and c in (a, b) for a, x, b, _, _ in lines)
# print(*sorted(c for a, x, b, _, c in lines if
#     x == "XOR" and all(d[0] not in 'xyz' for d in (a, b, c)) or
#     x == "AND" and not "x00" in (a, b) and r(c, 'XOR') or
#     x == "XOR" and not "x00" in (a, b) and r(c, 'OR') or
#     x != "XOR" and c[0] == 'z' and c != "z45"), sep=',')
