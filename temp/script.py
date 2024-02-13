import sys

def solution(str):
    return "dlrow"

if __name__ == "__main__":
    
    str = (sys.argv[0 + 1]) if len(sys.argv) > 0 + 1 else None
    
    output = (sys.argv[2]) if len(sys.argv) > 2 else None

    if solution(str) == output:
        exit(0)
    else:
        exit(1)