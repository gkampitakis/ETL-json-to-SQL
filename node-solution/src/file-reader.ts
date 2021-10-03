import { createReadStream, read } from 'fs';
import { createInterface } from 'readline';
import { Fn } from './types';
import Logger from './logger';

export function createStreamFileReader(path: string) {
  const stream = createReadStream(path, {
    encoding: 'utf-8',
    flags: 'r'
    // highWaterMark: 1,
    // <number> The maximum number of bytes to store in the internal buffer before ceasing to read from the underlying resource. Default: 16384 (16KB), or 16 for objectMode streams.
  });

  const reader = createInterface({
    input: stream,
    crlfDelay: Infinity,
    historySize: 0,
  });

  reader.on('pause', () => {
    console.log('reader is paused');
  });
  stream.on('pause', () => {
    console.log('stream is paused');
  });
  reader.on('close', () => Logger.info('Reader stream closing'));
  stream.on('end', () => Logger.info("Stream came to end"));

  stream.on('close', () => Logger.info('Stream stream closing'));
  stream.on('error', (error) => Logger.error(error));

  return {
    pause: () => reader.pause(),
    resume: () => reader.resume(),
    getData: (callback: (data: string) => void) => reader.on('line', callback),
    onEnd: (callback: Fn) => reader.on('close', callback)
  };
}