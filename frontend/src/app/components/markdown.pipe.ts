import {Pipe, PipeTransform} from '@angular/core';
import {marked, Tokens} from 'marked';

@Pipe({
  standalone: true,
  name: 'markdown'
})
export class MarkdownPipe implements PipeTransform {

  constructor() {
    let renderer = new marked.Renderer();
    renderer.link = function(l: Tokens.Link) {
      let link = marked.Renderer.prototype.link.call(this, l);
      return link.replace("<a","<a target='_blank' ");
    };

    marked.setOptions({
      renderer: renderer
    });
  }


  transform(value: string | undefined, ...args: unknown[]): unknown {
    if (value) {
      return marked.parse(value)
    } else {
      return "";
    }
  }

}
