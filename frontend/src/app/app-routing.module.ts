import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { SwipeComponent } from './components/swipe/swipe.component';
import { SelectCompanyComponent } from './components/select-company/select-company.component';
import { ReviewScreenComponent } from './components/review-screen/review-screen.component';
import { RegistrationCompleteComponent } from './components/registration-complete/registration-complete.component';
import { CompanyWelcomeComponent } from './components/company-welcome/company-welcome.component';
import { CompanyUniqueComponent } from './components/company-unique/company-unique.component';
import { AdminComponent } from './components/admin/admin.component';
import { CompanyEditorComponent } from './components/company-editor/company-editor.component';

const routes: Routes = [
  { path: '', component: SwipeComponent},
  { path: 'select', component: SelectCompanyComponent},
  { path: 'review', component: ReviewScreenComponent},
  { path: 'complete', component: RegistrationCompleteComponent},
  { path: 'employers', component: CompanyWelcomeComponent },
  { path: 'employers/:code', component: CompanyUniqueComponent },
  { path: 'admin', component: AdminComponent },
  { path: 'admin/companies', component: CompanyEditorComponent },
  { path: '**', redirectTo: '', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
