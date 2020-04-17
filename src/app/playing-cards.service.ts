import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders} from '@angular/common/http';
import { Deck, DeckGetResponse } from './deck';

@Injectable({
  providedIn: 'root'
})
export class PlayingCardsService {

  constructor(private http: HttpClient) { }

  getShuffledDeck(){

    return this.http.get<DeckGetResponse>('https://deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1')
  }

  shuffleDeck(deck_id: string){
    return this.http.get('https://deckofcardsapi.com/api/deck/' + deck_id + '/shuffle/')
  }

  drawCards(deck_id: string, num_cards: number){
    return this.http.get('https://deckofcardsapi.com/api/deck/' + deck_id + '/draw/?count=' + num_cards)
  }

}
 