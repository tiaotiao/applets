import sys
from scipy.sparse import csr_matrix
import numpy as np
from Eval import Eval
from math import log, exp
import time
from imdb import IMDBdata

class NaiveBayes:
    def __init__(self, data, ALPHA=1.0):
        self.ALPHA = ALPHA
        self.data = data # training data

        # Initalize parameters

        self.docs_len = data.X.shape[0]
        self.vocab_len = data.X.shape[1]  # k, total number of vocabulary

        self.deno_pos = 1.0
        self.deno_neg = 1.0

        self.P_word_positive = []
        self.P_word_negative = []

        self.Train(data.X, data.Y)

    # Train model - X are instances, Y are labels (+1 or -1)
    # X and Y are sparse matrices
    def Train(self, X, Y):
        # Estimate Naive Bayes model parameters
        num_docs = X.shape[0]
        num_vocab = X.shape[1]

        # Calculate P(+) and P(-)
        total_num_documents = 0.0
        self.num_positive_reviews = 0
        self.num_negative_reviews = 0

        for label in Y:
            total_num_documents += 1
            if label > 0:
                self.num_positive_reviews += 1
            else:
                self.num_negative_reviews += 1

        self.P_positive = self.num_positive_reviews / total_num_documents
        self.P_negative = self.num_negative_reviews / total_num_documents

        # print("TOTAL pos:", self.P_positive, self.num_positive_reviews, total_num_documents)
        # print("TOTAL neg:", self.P_negative, self.num_negative_reviews, total_num_documents)

        # Count words
        self.count_positive = np.zeros([1, num_vocab])
        self.count_negative = np.zeros([1, num_vocab])

        for doc_idx in range(num_docs):
            label = Y[doc_idx]
            row = X.getrow(doc_idx)    # row[wordId] = count
            if label > 0:
                self.count_positive = self.count_positive + row
            else:
                self.count_negative = self.count_negative + row

        # Count total words in each category
        self.total_positive_words = self.count_positive.sum()
        self.total_negative_words = self.count_negative.sum()

        # Calculate P(w|+) and P(w|-)
        self.P_word_positive = []
        self.P_word_negative = []
        sum_pos = 0
        sum_neg = 0

        for wordId in range(self.vocab_len):
            count_pos = self.count_positive[0, wordId]
            count_neg = self.count_negative[0, wordId]

            P_word_pos = (count_pos + self.ALPHA) / (self.total_positive_words + self.vocab_len * self.ALPHA)
            P_word_neg = (count_neg + self.ALPHA) / (self.total_negative_words + self.vocab_len * self.ALPHA)

            sum_pos += P_word_pos
            sum_neg += P_word_neg

            self.P_word_positive.append(P_word_pos)
            self.P_word_negative.append(P_word_neg)

        # Normalize
        for i in range(self.vocab_len):
            self.P_word_positive[i] /= sum_pos
            self.P_word_negative[i] /= sum_neg

        # debug
        # sum_pos = np.sum(self.P_word_positive)
        # sum_neg = np.sum(self.P_word_negative)
        # print("Train:", sum_pos, sum_neg, self.P_positive, self.P_negative)

        return

    # Predict labels for instances X
    # Return: Sparse matrix Y with predicted labels (+1 or -1)
    def PredictLabel(self, X, probThresh=0.5):
        num_docs = X.shape[0]

        pred_probs = self.PredictProb(X, [i for i in range(num_docs)])
        
        pred_labels = []

        for i in range(num_docs):
            prob = pred_probs[i]
            if prob > probThresh:
                pred_labels.append(1.0)
            else:
                pred_labels.append(-1.0)

        return pred_labels

        
    # def LogSum(self, log_list):   
        # Return log(x+y), avoiding numerical underflow/overflow.
        # https://stats.stackexchange.com/questions/105602/example-of-how-the-log-sum-%20exp-trick-works-in-naive-bayes
        # n = len(log_list)
        # m = log_list[0]
        # for v in log_list:
        #     m = max(m, v)
        # expSum = 0.0
        # for v in log_list:
        #     expSum += exp(v - m)
        # return m + log(expSum)

    def LogSum(self, logx, logy):   
        m = max(logx, logy)        
        return m + log(exp(logx - m) + exp(logy - m))

    # Predict the probability of each indexed review in sparse matrix text
    # of being positive
    # Prints results
    def PredictProb(self, X, indexes):
        pred_probs = []

        for i in indexes:
            # TO DO: Predict the probability of the i_th review in test being positive review
            # TO DO: Use the LogSum function to avoid underflow/overflow
            doc = X[i]
            indices = doc.nonzero()

            # Init numerator
            nume_pos = log(self.P_positive)
            nume_neg = log(self.P_negative)

            for idx in range(len(indices[0])):
                wordId = indices[1][idx]
                count = doc[0, wordId]
                # Get probability of word
                P_word_pos = self.P_word_positive[wordId]
                P_word_neg = self.P_word_negative[wordId]
                # Sum up numerator
                nume_pos += log(P_word_pos) * count
                nume_neg += log(P_word_neg) * count
                # print("    >>> ", P_word_pos, P_word_neg, count)

            # Calc denominator using LogSum
            deno = self.LogSum(nume_pos, nume_neg)
            
            predicted_prob_positive = exp(nume_pos - deno)
            predicted_prob_negative = exp(nume_neg - deno)
            
            # if predicted_prob_positive > predicted_prob_negative:
            #     predicted_label = 1.0
            # else:
            #     predicted_label = -1.0
            
            # print("predict prob", predicted_prob_positive, predicted_prob_negative, nume_pos, nume_neg, deno)
            
            # add result
            pred_probs.append(predicted_prob_positive)

            # print test.Y[i], test.X_reviews[i]
            # TODO: Comment the line above, and uncomment the line below
            # print(test.Y[i], predicted_label, predicted_prob_positive, predicted_prob_negative)
        
        return pred_probs

    def EvalCount(self, Y_pred, Y):
        truePos = 0.0
        trueNeg = 0.0
        falsePos = 0.0
        falseNeg = 0.0
        for i in range(len(Y_pred)):
            pred = Y_pred[i]
            label = Y[i]
            if pred == label:
                if pred == True:
                    truePos += 1
                else:
                    trueNeg += 1
            else:
                if pred == True:
                    falsePos += 1
                else:
                    falseNeg += 1

        # print("EvalCount", truePos, trueNeg, falsePos, falseNeg)
        return truePos, trueNeg, falsePos, falseNeg        


    def EvalPrecision(self, Y_pred, test):
        # What percent of positive predictions were correct? 
        # TP / (TP + FP)
        tp, tn, fp, fn = self.EvalCount(Y_pred, test.Y)
        deno = tp + fp
        if deno <= 0:
            return -1
        return tp / deno


    def EvalRecall(self, Y_pred, test):
        # What percent of the positive cases did you catch? 
        # TP / (TP + FN)
        tp, tn, fp, fn = self.EvalCount(Y_pred, test.Y)
        deno = tp + fn
        if deno <= 0:
            return -1
        return tp / deno

    # Evaluate performance on test data 
    def Eval(self, test, probThresh=0.5):
        Y_pred = self.PredictLabel(test.X, probThresh)
        ev = Eval(Y_pred, test.Y)

        prec = self.EvalPrecision(Y_pred, test)
        print("Test Precision: ", prec)

        recall = self.EvalRecall(Y_pred, test)
        print("Test Recall: ", recall)

        acc = ev.Accuracy()
        print("Test Accuracy: ", acc)

    def PrintTopWords(self, prob_list1, prob_list2, top=20):
        # reorganize data
        sort_list = []
        for wordId in range(self.vocab_len):
            weight = log(prob_list1[wordId]) - log(prob_list2[wordId])
            sort_list.append((weight, wordId))
        # sort
        sort_list = sorted(sort_list, key=lambda x: x[0], reverse=True)[:top]
        # print
        # print(sort_list)
        print_list = []
        for i in range(top):
            weight = sort_list[i][0]
            wordId = sort_list[i][1]
            vocab = self.data.vocab.GetWord(wordId)
            print_list.append(vocab)
            print_list.append(weight)
            # print("   %d:\t%15s %d\t- %f" % (i, vocab, wordId, weight))
        print(print_list)

    def PrintProbReviews(self, test, count=10):
        pred_probs = self.PredictProb(test.X, [i for i in range(count)])
        for i in range(count):
            prob = pred_probs[i]
            label = test.Y[i]
            print("Predict Probability: ", i, prob, label)


if __name__ == "__main__":
    
    path = "data/aclImdb"
    ALPHA = 1.0

    if len(sys.argv) > 1:
        path = sys.argv[1]
    if len(sys.argv) > 2:
        ALPHA = float(sys.argv[2])

    print("Reading Training Data")
    traindata = IMDBdata("%s/train" % path)
    print("Reading Test Data")
    testdata  = IMDBdata("%s/test" % path, vocab=traindata.vocab)    
    print("Computing Parameters")
    nb = NaiveBayes(traindata, ALPHA)
    
    # probThresh_list = [0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9]
    probThresh_list = [0.5]
    for probThresh in probThresh_list:
        print("Evaluating", probThresh)
        nb.Eval(testdata, probThresh)

    print("Predict Probability")
    nb.PrintProbReviews(testdata)

    top = 20
    print("Top %d positive words" % (top))
    nb.PrintTopWords(nb.P_word_positive, nb.P_word_negative, top)
    print("Top %d negative words" % (top))
    nb.PrintTopWords(nb.P_word_negative, nb.P_word_positive, top)

