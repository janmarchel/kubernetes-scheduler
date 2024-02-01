import unittest
from utils import create_sequences

class TestUtils(unittest.TestCase):
    def test_create_sequences(self):
        data = [i for i in range(10)]
        X, y = create_sequences(data, 5)
        self.assertEqual(X.shape[0], y.shape[0])

if __name__ == '__main__':
    unittest.main()