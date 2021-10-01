import { createReadStream } from 'fs';
import { createInterface } from 'readline';
import { Fn } from './types';

export function createStreamFileReader(path: string) {
  const stream = createReadStream(path, {
    encoding: 'utf-8',
    flags: 'r',
    highWaterMark: 1
    // <number> The maximum number of bytes to store in the internal buffer before ceasing to read from the underlying resource. Default: 16384 (16KB), or 16 for objectMode streams.
  });

  const reader = createInterface({
    input: stream,
    crlfDelay: Infinity
  });

  return {
    pause: () => reader.pause(),
    resume: () => reader.resume(),
    getData: (callback: (data: string) => void) => reader.on('line', callback),
    onEnd: (callback: Fn) => stream.on('end', callback)
  };
}