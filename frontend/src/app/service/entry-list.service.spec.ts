import { TestBed } from '@angular/core/testing';

import { EntryListServiceService } from './entry-list-service.service';

describe('EntryListServiceService', () => {
  let service: EntryListServiceService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EntryListServiceService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
