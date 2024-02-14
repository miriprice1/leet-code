const args = process.argv.slice(2);

function solution(arr){
  return arr[0]  
};



const arr = JSON.parse(args[0])



const output = parseFloat(args[1])

if (solution(arr) == output){
    process.exit(0)
}
else{
    process.exit(1)
}