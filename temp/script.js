const args = process.argv.slice(2);

function solution(arr){
  return arr[arr.length-1]  
};




//case of array string or boolean
let tmp_arr = args[0].slice(1 ,-1);
let arr = tmp_arr.split(',').map(elem => elem.trim());









//case of string
let output = args[1]






if (solution(arr) == output){
    process.exit(0)
}
else{
    process.exit(1)
}