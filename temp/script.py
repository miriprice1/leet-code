import sys
import json

def solution(number1,number2):
    return number1 + number2

if __name__ == "__main__":
    args = sys.argv[1:]
    
    
    number1 = float(args[0])
    
    
    
    number2 = float(args[1])
    
    
    
    output = float(args[2])
    

    if solution(number1, number2) == output:
        exit(0)
    else:
        exit(1)