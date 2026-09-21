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
package ru.vorobevaqa.controller.api;

import java.util.List;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import ru.vorobevaqa.dto.SkillDto;
import ru.vorobevaqa.entity.SkillCategory;
import ru.vorobevaqa.service.SkillService;

@RestController
@RequestMapping("/api/skills")
@RequiredArgsConstructor
public class SkillRestController {

  private final SkillService skillService;

  @GetMapping
  public ResponseEntity<List<SkillDto>> getSkills(
      @RequestParam(value = "category", required = false) SkillCategory category) {
    if (category != null) {
      return ResponseEntity.ok(skillService.getSkillsByCategory(category));
    }
    return ResponseEntity.ok(skillService.getAllSkills());
  }
}
