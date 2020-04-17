import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { PlayingCardsService } from '../playing-cards.service';
import { Deck, DeckGetResponse } from '../deck';
import {CdkDragDrop, moveItemInArray} from '@angular/cdk/drag-drop';

interface Card{
  deck_id: string,
  value: string,
  suit: string,
  code: string,
  image: string
}

@Component({
  selector: 'app-game',
  templateUrl: './game-room.component.html',
  styleUrls: ['./game-room.component.sass']
})
export class GameRoomComponent implements OnInit {

  deckInt: DeckGetResponse;
  deckObj: Deck = new Deck(true, "viwzpbuz4n22", true, 52, [], null);

  constructor(private http: HttpClient, private playCardServ: PlayingCardsService) { }

  cardList: Card[] = [];

  ngOnInit(): void {
    
  }

  initShuffledDeck(){
    this.playCardServ.getShuffledDeck().subscribe(res => {

      console.log(res)

      this.deckObj = new Deck(res.success, res.deck_id, res.shuffled, res.remaining, [], null) 
    })
  }

  shuffleDeck(){
    this.playCardServ.shuffleDeck(this.deckObj.deck_id).subscribe(res => {
      
      console.log(res)
    })
  }

  drawCards(){
    console.log("ping")
    this.playCardServ.drawCards(this.deckObj.deck_id, 2).subscribe(res => {
      res['cards'].forEach(card => {
        let drawnCard = <Card>{
          deck_id: res['deck_id'],
          value: card['value'],
          suit: card['suit'],
          code: card['code'],
          image: card['image']
        }

        this.cardList.push(drawnCard)

      });

      console.log(this.cardList)
    })
  }

  onTaskDrop(event: CdkDragDrop<string[]>){
    moveItemInArray(this.cardList, event.previousIndex, event.currentIndex)
  }

}
