import sys
import json

def solution(arr):
    return arr[0]

if __name__ == "__main__":
    args = sys.argv[1:]
      
    arr = json.loads(args[0])  
    output = float(args[1])
    
    if solution(arr) == output:
        exit(0)
    else:
        exit(1)