const args = process.argv.slice(2);

function solution(number1,number2){
  return number1 + number2  
};



const number1 = parseFloat(args[0])



const number2 = parseFloat(args[1])



const output = parseFloat(args[2])

if (solution(number1, number2) == output){
    process.exit(0)
}
else{
    process.exit(1)
}