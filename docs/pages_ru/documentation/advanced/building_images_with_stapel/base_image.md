---
title: Базовый образ
sidebar: documentation
permalink: documentation/advanced/building_images_with_stapel/base_image.html
author: Alexey Igrychev <alexey.igrychev@flant.com>
summary: |
  <a class="google-drawings" href="../../../images/configuration/base_image1.png" data-featherlight="image">
      <img src="../../../images/configuration/base_image1_preview.png">
  </a>

  <div class="language-yaml highlighter-rouge"><div class="highlight"><pre class="highlight"><code><span class="na">from</span><span class="pi">:</span> <span class="s">&lt;image[:&lt;tag&gt;]&gt;</span>
  <span class="na">fromLatest</span><span class="pi">:</span> <span class="s">&lt;bool&gt;</span>
  <span class="na">fromCacheVersion</span><span class="pi">:</span> <span class="s">&lt;arbitrary string&gt;</span>
  <span class="na">fromImage</span><span class="pi">:</span> <span class="s">&lt;image name&gt;</span>
  <span class="na">fromArtifact</span><span class="pi">:</span> <span class="s">&lt;artifact name&gt;</span>
  </code></pre></div>
  </div>
---

Пример минимального `werf.yaml`:
```yaml
project: my-project
configVersion: 1
---
image: example
from: alpine
```

Приведенная конфигурация описывает _образ_ `example`, _базовым образом_ для которого является образ с именем `alpine`.

_Базовый образ_ может быть указан с помощью директив `from`, `fromImage` или `fromArtifact`.

## from, fromLatest

Директива `from` определяет имя и тег _базового образа_. Если тег не указан, то по умолчанию — `latest`.

```yaml
from: <image>[:<tag>]
```

По умолчанию процесс сборки не зависит от digest'а _базового образа_, а зависит только от значения директивы `from`.
Поэтому изменение _базового образа_ в локальном хранилище или в Docker registry не будет влиять на сборку, пока стадия _from_, с указанным значением образа, находится в _stages storage_.

Если вам нужна проверка digest образа, чтобы всегда использовать актуальный _базовый образ_, вы можете использовать директиву `fromLatest`.
Это приведет к тому, что при каждом запуске werf будет проверяться актуальный digest _базового образа_ в Docker registry.

Пример использования директивы `fromLatest`:

```yaml
fromLatest: true
```

> Обратите внимание, что если вы включаете _fromLatest_, то werf начинает использовать digest актуального _базового образа_ при подсчете дайджеста стадии _from_.
> Это может приводить к неконтролируемым сменам дайджестов стадий: все образы стадий, собранные ранее, становятся неактуальными, если меняется базовый образ в репозитории.
> Примеры проблем, которые может вызвать это поведение в CI процессах (например, в pipeline GitLab):
>
> * Сборка прошла успешно, но затем обновляется _базовый образ_, и **следующие задания pipeline** (например, деплой) уже не работают. Это происходит потому, что еще не существует конечного образа, собранного с учетом обновленного _базового образа_.
> * Собранное приложение успешно развернуто, но затем обновляется _базовый образ_, и **повторный запуск** деплоя уже не работает. Это также происходит потому, что еще не существует конечного образа, собранного с учетом обновленного _базового образа_.
>
> Если вы всё же хотите использовать функционал данной директивы, то необходимо выключить режим детерминизма в werf с помощью опции --disable-determinism.
>
> **Крайне не рекомендуется использовать актуальный базовый образ таким способом**. Используйте конкретный неизменный tag или периодически обновляйте значение [fromCacheVersion](#fromcacheversion) для обеспечения предсказуемого и контролируемого жизненного цикла приложения

## fromImage и fromArtifact

В качестве _базового образа_ можно указывать не только образ из локального хранилища или Docker registry, но и имя другого _образа_ или [_артефакта_]({{ "documentation/advanced/building_images_with_stapel/artifacts.html" | relative_url }}), описанного в том же файле `werf.yaml`. В этом случае необходимо использовать директивы `fromImage` и `fromArtifact` соответственно.

```yaml
fromImage: <image name>
fromArtifact: <artifact name>
```

Если _базовый образ_ уникален для конкретного приложения, то рекомендуемый способ — хранить его описание в конфигурации приложения (в файле `werf.yaml`) как отдельный _образ_ или _артефакт_, вместо того, чтобы ссылаться на Docker-образ.

Также эта рекомендация будет полезной, если вам, по каким-либо причинам, не хватает существующего _конвейера стадий_.
Используя в качестве _базового образа_ образ, описанный в том же `werf.yaml`, вы по сути можете построить свой _конвейер стадий_.

<a class="google-drawings" href="{{ "images/configuration/base_image2.png" | relative_url }}" data-featherlight="image">
<img src="{{ "images/configuration/base_image2_preview.png" | relative_url }}">
</a>

## fromCacheVersion

Как описано выше, в обычном случае процесс сборки активно использует кэширование.
При сборке выполняется проверка — изменился ли _базовый образ_.
В зависимости от используемых директив эта проверка на изменение digest или имени и тега образа.
Если образ не изменился, то дайджест стадии `from` остается прежней, и если в _stages storage_ есть образ с таким дайджестом, то он и будет использован при сборке.

С помощью директивы `fromCacheVersion` вы можете влиять на дайджест стадии `from` (т.к. значение `fromCacheVersion` — это часть дайджеста стадии) и, таким образом, управлять принудительной пересборкой образа.
Если вы измените значение, указанное в директиве `fromCacheVersion`, то независимо от того, менялся _базовый образ_ (или его digest) или остался прежним, при сборке изменится дайджест стадии `from` и, соответственно, всех последующих стадий.
Это приведет к тому, что сборка всех стадий будет выполнена повторно.

```yaml
fromCacheVersion: <arbitrary string>
```
