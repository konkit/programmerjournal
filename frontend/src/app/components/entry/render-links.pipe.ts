import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  standalone: true,
  name: 'renderLinks'
})
export class RenderLinksPipe implements PipeTransform {

  transform(value: string, ...args: unknown[]): string {
    return value
      .split(" ")
      .flatMap(str => {
        if (str.startsWith("http://") || str.startsWith("https://")) {
          return `<a href="${str}" target="_blank" rel="noopener noreferrer">${str}</a>`
        } else if (str.startsWith("#")) {
          return `<a href="tags/${str.substring(1).toLowerCase()}">${str}</a>`
        } else {
          return str
        }
      })
      .join(" ")
  }

}
