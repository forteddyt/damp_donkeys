import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SwipeComponent } from './components/swipe/swipe.component';
import { SelectCompanyTileComponent } from './components/select-company-tile/select-company-tile.component';
import { SelectCompanyComponent } from './components/select-company/select-company.component';
import { ReviewScreenComponent } from './components/review-screen/review-screen.component';
import { RegistrationCompleteComponent } from './components/registration-complete/registration-complete.component';
import { CompanyUniqueComponent } from './components/company-unique/company-unique.component';
import { CompanyWelcomeComponent } from './components/company-welcome/company-welcome.component';
import { AdminComponent } from './components/admin/admin.component';
import { AdminWelcomeComponent } from './components/admin-welcome/admin-welcome.component';
import { CompanyEditorComponent } from './components/company-editor/company-editor.component';
import { StatisticsComponent } from './components/statistics/statistics.component';
import { AddCompanyComponent } from './components/add-company/add-company.component';
import bootstrap from "bootstrap";
import { InterviewComponent } from './components/interview/interview.component';
import { CompanyNameTileComponent } from './components/company-name-tile/company-name-tile.component';

@NgModule({
  declarations: [
    AppComponent,
    SwipeComponent,
    SelectCompanyComponent,
    ReviewScreenComponent,
    RegistrationCompleteComponent,
    CompanyUniqueComponent,
    CompanyWelcomeComponent,
    AdminComponent,
    AdminWelcomeComponent,
    CompanyEditorComponent,
    StatisticsComponent,
    AddCompanyComponent,
    SelectCompanyTileComponent,
    InterviewComponent,
    CompanyNameTileComponent
  ],
  entryComponents: [
    SelectCompanyTileComponent,
    InterviewComponent,
    CompanyNameTileComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
