require_relative './data_mixin'

class Repo
  include Virtus.model
  include DataMixin

  def record
    connection.exec("INSERT INTO repos (id, data) VALUES ($1, $2)", [id, data])
  end

  def url
    data_hash['url']
  end

  def is_dup?
    result = connection.exec("SELECT id FROM repos WHERE id = $1", [id])
    result.ntuples == 1
  end
end