/*
 * Copyright 2024-2026 Andrey Vorobev
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
package ru.vorobevaqa.qa.ui;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import ru.vorobevaqa.qa.pages.EducationPage;

@Epic("Web Application")
@Feature("Education & Certificates")
@Story("Просмотр образования и дипломов")
@Tag("ui")
public class EducationUiTest extends BaseUiTest {

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-EDU-01: Отображение страницы образования и списка дипломов")
  @Description(
      "Проверка рендеринга заголовка квалификации и карточек дипломов государственного образца")
  void shouldDisplayEducationPageAndDiplomas() {
    new EducationPage().openPage().shouldSeeEducationHeader().shouldHaveDiplomaCards();
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-EDU-02: Открытие модального окна просмотра диплома")
  @Description("Проверка открытия увеличенного изображения диплома в модальном окне при клике")
  void shouldOpenCertificateModalOnCardClick() {
    new EducationPage().openPage().openFirstCertificateModal();
  }
}
