import sys

def solution(number1,number2):
    return number1 + number2

if __name__ == "__main__":
    
    number1 = float(sys.argv[0 + 1]) if len(sys.argv) > 0 + 1 else None
    
    number2 = float(sys.argv[1 + 1]) if len(sys.argv) > 1 + 1 else None
    
    output = float(sys.argv[3]) if len(sys.argv) > 3 else None

    if solution(number1, number2) == output:
        exit(0)
    else:
        exit(1)