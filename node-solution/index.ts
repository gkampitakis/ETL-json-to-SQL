import Logger from './src/logger';
import { readFileSync, createReadStream } from 'fs';
import { createInterface } from 'readline';


Logger.info('Hello world');

// const file = readFileSync('../data-to-load/yelp_academic_dataset_review.json', {
//   encoding: 'utf-8'
// });


const stream = createReadStream('../data-to-load/yelp_academic_dataset_review.json', {
  encoding: 'utf-8',
  flags: 'r'
});


async function test() {
  let count = 0;
  const lineReader = createInterface({
    input: stream,
    crlfDelay: Infinity
  });


  for await (const line of lineReader) {
    console.log(JSON.parse(line)["review_id"]);
    count++;
  }

  console.log(count);
}

test();




// stream.on('data', (data) => {
//   tmp++;
//   console.log(data);

//   if (tmp > 1) {
//     stream.close();
//   }
// });
// stream.on('error', console.error);
// stream.on('close', () => {
//   console.log('stream closing');
// });


// Load files , gradually
// Transform them
// Push them