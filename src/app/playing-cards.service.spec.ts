import { TestBed } from '@angular/core/testing';

import { PlayingCardsService } from './playing-cards.service';

describe('PlayingCardsService', () => {
  let service: PlayingCardsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PlayingCardsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
