const fs = require('fs');
const path = require('path');

function replaceInFile(filePath) {
  let content = fs.readFileSync(filePath, 'utf8');
  const target = '333f6432-a01c-412f-99f4-0f08ca0d8eb1';
  const replacement = '99ef0062-1040-4574-a4be-0123abce5670';
  if (content.includes(target)) {
    content = content.split(target).join(replacement);
    fs.writeFileSync(filePath, content, 'utf8');
    console.log(`Replaced in ${filePath}`);
  }
}

const files = [
  'd:\\chia.florist\\e-commerce\\app\\composables\\useCart.ts',
  'd:\\chia.florist\\e-commerce\\app\\pages\\checkout.vue',
  'd:\\chia.florist\\e-commerce\\app\\pages\\products\\custom.vue',
  'd:\\chia.florist\\e-commerce\\app\\pages\\products\\[id].vue'
];

files.forEach(f => replaceInFile(f));
console.log('Done replacing.');
