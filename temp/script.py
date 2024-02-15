import sys
import json

def solution(arr):
    return arr[len(arr)-1]

if __name__ == "__main__":
    args = sys.argv[1:]
    
    
    
    tmp_arr = args[0][0+1:-1]
    arr = [elem.strip() for elem in tmp_arr.split(',')]
    
    
    

    
    
    output = args[1]
    
    

    if solution(arr) == output:
        exit(0)
    else:
        exit(1)


