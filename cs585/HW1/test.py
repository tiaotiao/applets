
from scipy.sparse import csr_matrix
import numpy as np

def testCSRMatrix():
    pass
    row = [0,0,0,0, 1,1,1,1]
    col = [0,1,0,1, 2,1,2,1]
    data = [1,5,1,0, 2,0,5,0]

    X = csr_matrix((data, (row, col)), shape=(max(row)+1, 3))

    print(X)

    line = X[0]

    print("-----")
    print(line)

    print("-----")
    print(line.nonzero())

    s = np.zeros([1, X.shape[1]])

    print("=====")
    print(s)
    for i in range(X.shape[0]):
        row = X.getrow(i)
        s = s + row
        print(">>>>>>>>> row ", i, "\n", row, s)
        print("+++", row[0][0])

    print("sum=", s.sum())
        

def main():
    testCSRMatrix()

if __name__ == "__main__":
    main()