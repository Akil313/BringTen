import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { GameRoomComponent } from './game/game-room.component';


const routes: Routes = [
  {
    path: '', component: HomeComponent
  },
  {
    path: 'game', component: GameRoomComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
