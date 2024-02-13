function solution(number1,number2){
  return number1 + number2;  
};


const number1 = parseFloat(process.argv[0 + 2])

const number2 = parseFloat(process.argv[1 + 2])

const output = parseFloat(process.argv[2 + 2])

if (solution(number1, number2) == output){
    process.exit(0)
}
else{
    process.exit(1)
}