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
package ru.vorobevaqa.dto;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import ru.vorobevaqa.entity.Skill;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SkillDto {

  private Long id;
  private String name;
  private String categoryKey;
  private String categoryDisplayName;
  private String levelTitle;
  private int levelPercentage;
  private String badgeColor;
  private String description;
  private String iconClass;
  private List<String> tools;

  public static SkillDto fromEntity(Skill skill) {
    List<String> toolsList =
        (skill.getKeyTools() != null && !skill.getKeyTools().isBlank())
            ? Arrays.stream(skill.getKeyTools().split(","))
                .map(String::trim)
                .filter(s -> !s.isEmpty())
                .toList()
            : Collections.emptyList();

    return SkillDto.builder()
        .id(skill.getId())
        .name(skill.getName())
        .categoryKey(skill.getCategory().name())
        .categoryDisplayName(skill.getCategory().getDisplayName())
        .levelTitle(skill.getLevel().getTitle())
        .levelPercentage(skill.getLevel().getPercentage())
        .badgeColor(skill.getLevel().getBadgeColor())
        .description(skill.getDescription())
        .iconClass(skill.getIconClass())
        .tools(toolsList)
        .build();
  }
}
