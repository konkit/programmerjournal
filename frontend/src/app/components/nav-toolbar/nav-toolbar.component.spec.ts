import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NavToolbarComponent } from './nav-toolbar.component';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

describe('ToolbarComponent', () => {
  let component: NavToolbarComponent;
  let fixture: ComponentFixture<NavToolbarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NavToolbarComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({ id: 'test-id' })
          }
        }
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NavToolbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('renders a link for every section', () => {
    const links = Array.from(fixture.nativeElement.querySelectorAll('.nav-links a') as NodeListOf<HTMLAnchorElement>);

    expect(links.map(a => a.getAttribute('href'))).toEqual(component.links.map(l => l.path));
    expect(links.map(a => a.textContent!.trim())).toEqual(component.links.map(l => l.label));
  });

  it('has a labelled menu button for narrow screens', () => {
    const button = fixture.nativeElement.querySelector('.nav-menu-button') as HTMLButtonElement;
    expect(button.getAttribute('aria-label')).toBe('Open navigation menu');
  });
});
