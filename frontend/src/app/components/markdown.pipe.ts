import {Pipe, PipeTransform} from '@angular/core';
import {marked} from 'marked';

@Pipe({
  standalone: true,
  name: 'markdown'
})
export class MarkdownPipe implements PipeTransform {

  transform(value: string | undefined, ...args: unknown[]): unknown {
    if (value) {
      return marked.parse(value)
    } else {
      return "";
    }
  }

}
